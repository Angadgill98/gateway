import React, { useEffect, useState } from 'react'
import { useGlobal } from '../../context/GLobal';


interface Dashboard_props {

}


const Dashboard = ({}: Dashboard_props) => {
    const { app } = useGlobal();

    const [isEditing, setIsEditing] = useState(false);
    const [rateLimit, setRateLimit] = useState(0);
    const [refillRate, setRefillRate] = useState(0);

    useEffect(() => {
        if (app) {
            setRateLimit(app.rate_limit);
            setRefillRate(app.refill_rate);
        }
        setIsEditing(false);
    }, [app]);

    if (!app) {
        return <div>No app selected</div>;
    }

    const handleCancel = () => {
        setRateLimit(app.rate_limit);
        setRefillRate(app.refill_rate);
        setIsEditing(false);
    };

    return (
        <div>
            <h2>{app.name}</h2>

            {isEditing ? (
                <div>
                    <div>
                        <label>Rate Limit</label>
                        <input
                            type="number"
                            value={rateLimit}
                            onChange={(e) => setRateLimit(Number(e.target.value))}
                        />
                    </div>

                    <div>
                        <label>Refill Rate</label>
                        <input
                            type="number"
                            value={refillRate}
                            onChange={(e) => setRefillRate(Number(e.target.value))}
                        />
                    </div>

                    <button onClick={() => setIsEditing(false)}>Save</button>
                    <button onClick={handleCancel}>Cancel</button>
                </div>
            ) : (
                <div>
                    <div>Rate Limit: {app.rate_limit}</div>
                    <div>Refill Rate: {app.refill_rate}</div>

                    <button onClick={() => setIsEditing(true)}>Edit</button>
                </div>
            )}
        </div>
    );
};

export default Dashboard;