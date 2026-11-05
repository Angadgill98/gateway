import React, { useEffect, useRef } from 'react'
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


const EndpointTemplate = ({ endpoint }: { endpoint: Endpoint }) => {
    return (
        <li key={endpoint.id}>
            <div>{endpoint.route}</div>
            <div>Rate Limit: {endpoint.rate_limit}</div>
            <div>Refill Rate: {endpoint.refill_rate}</div>
        </li>
    );
};