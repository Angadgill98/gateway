import React, { useEffect, useRef, useState } from 'react'
import { useGlobal } from '../../context/GLobal';
import api from '../../Api/Api';
import type { Endpoint } from '../../models';


interface Dashboard_props {

}


const Dashboard = ({}: Dashboard_props) => {
    const { app, endpoints, setEndpoints } = useGlobal();
    const endpointRef = useRef<HTMLInputElement>(null);

    useEffect(() => {
        if (!app) {
            setEndpoints([]);
            return;
        }

        const fetchEndpoints = async () => {
            const result = await api.GetAllEndpoints(app.id);

            if (!result.success) {
                console.error(result.message);
                return;
            }

            setEndpoints(result.data ?? []);
        };

        fetchEndpoints();
    }, [app?.id]);

    if (!app) {
        return <div>No app selected</div>;
    }

    

    const handleCreate_endpoint = async () => {
        const endpoint = endpointRef.current?.value.trim();

        if (!endpoint) {
            return;
        }

        const result = await api.CreateEndpoint(app.id, endpoint);

        if (!result.success) {
            console.error(result.message);
            return;
        }
        setEndpoints(prev => [
            ...prev,
            {
                id: result.data.id,
                app_id: app.id,
                route: endpoint,
                destination_ids: [],
                rate_limit: 100,
                refill_rate: 10
            }
        ]);

        if (endpointRef.current) {
            endpointRef.current.value = "";
        }
    };
    return (
        <div>
            <h2>{app.name}</h2>
            <div>Rate Limit: {app.rate_limit}</div>
            <div>Refill Rate: {app.refill_rate}</div>




            <input ref={endpointRef} type='text' placeholder='Enter Endpoint'></input>
            <button onClick={handleCreate_endpoint}>Create Endpoint</button>

            {endpoints.map((endpoint)=>(<EndpointTemplate key={endpoint.id} endpoint={endpoint} />))}
        </div>
    );
};

export default Dashboard;

export interface Destination {
    id: string;
    endpoint_id: string;
    url: string;
    weight: number;
    is_active: boolean;
}
const EndpointTemplate = ({ endpoint }: { endpoint: Endpoint }) => {
    let [show,setshow]=useState<boolean>(false)
    let destinationRef = useRef<HTMLInputElement>(null);
    let weightRef = useRef<HTMLInputElement>(null);
    const [destinations, setDestinations] = useState<Destination[]>([]);
    useEffect(() => {
        if (!show) return;

        const fetchDestinations = async () => {
            const result = await api.GetEndpointDestinations(endpoint.id);

            if (!result.success) {
                console.error(result.message);
                return;
            }

            setDestinations(result.data ?? []);
        };

        fetchDestinations();
    }, [show, endpoint.id]);

    let handleSave = async () => {
        const url = destinationRef.current?.value.trim();
        const weight = Number(weightRef.current?.value);

        if (!url || !weight || weight <= 0) return;

        const result = await api.CreateDestination(endpoint.id, url, weight);

        if (!result.success) {
            console.error(result.message);
            return;
        }

        if (destinationRef.current) destinationRef.current.value = "";
        if (weightRef.current) weightRef.current.value = "";
    };


    return (
        <li key={endpoint.id}>
            <div onClick={() => setshow(prev => !prev)}>
                <div>{endpoint.route}</div>
                <div>Rate Limit: {endpoint.rate_limit}</div>
                <div>Refill Rate: {endpoint.refill_rate}</div>
            </div>

            {show && 
            <div>
                <input ref={destinationRef} type="text" placeholder="Enter Destination URL" />
                <input ref={weightRef} type="number" placeholder="Enter Weight" />
                <button onClick={handleSave}>Save</button>

                <ul>
                    {destinations.map(destination => (
                        <DestinationTemplate key={destination.id} destination={destination} endpoint={endpoint.route} />
                    ))}
                </ul>
            </div>}
        </li>
    );
};


const DestinationTemplate = ({ destination ,endpoint}: { destination: Destination,endpoint:string }) => {
    let [Destination,setdestination]=useState<Destination>(destination)

    let {app}=useGlobal()
    let handleToggleActive = async () => {
        const updatedActive = !Destination.is_active;

         const result = await api.UpdateDestination(
            Destination.id,
            Destination.endpoint_id,
            app!.name,
            endpoint,
            Destination.url,
            Destination.weight,
            updatedActive
        );

        if (!result.success) {
            console.error(result.message);
            return;
        }

        setdestination(prev => ({
            ...prev,
            is_active: !prev.is_active
        }));
    };
    return (
        <li>
            <div>URL: {Destination.url}</div>
            <div>Weight: {Destination.weight}</div>
            <div>Active: {Destination.is_active ? "Yes" : "No"}</div>

             <button onClick={handleToggleActive}>
                {Destination.is_active ? "Deactivate" : "Activate"}
            </button>
        </li>
    );
};