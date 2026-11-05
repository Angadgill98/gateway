import React, { useEffect, useRef, useState } from 'react'
import type { App } from '../../models'
import api from '../../Api/Api'
import { useGlobal } from '../../context/GLobal'


interface Sidebar_props {
    classes:any
}

const Sidebar = ({classes} :Sidebar_props) => {
    let [apps,setApps]=useState<App[]>([])
    const appNameRef = useRef<HTMLInputElement>(null);

    let handleCreateApp=async ()=>{
        let name = appNameRef.current?.value.trim();

        if (!name) return;

        if (appNameRef.current) {
            appNameRef.current.value = "";
        }

        let res=await api.RegisterApp(name)
        if (!res.success || !res.data) return;

        const newApp: App = {
            id: res.data.id,
            name: name,
            rate_limit: 100,
            refill_rate: 10
        };

        setApps(prev => [...prev, newApp]);

        if (appNameRef.current) {
            appNameRef.current.value = "";
        }
    }

    useEffect(() => {
        const fetchApps = async () => {
            try {
                const res = await api.GetAllApps();

                if (!res.success || !res.data) {
                    setApps([]);
                    return;
                }

                setApps(res.data);
            } catch (error) {
                console.error("Failed to fetch apps:", error);
                setApps([]);
            }
        };

        fetchApps();
    }, []);

  return (
    <div className={classes.sidebar_container}>
        <input ref={appNameRef} type='text' placeholder='Enter App Name'></input>
        <button onClick={handleCreateApp}>Add</button>
        
        <ul>
            {apps.map(app => <AppTemplate key={app.id} app={app} />)}
        </ul>
    </div>
  )
}

export default Sidebar


const AppTemplate = ({app}:{app:App})=>{

    const { setApp } = useGlobal();

    let set_as_currrent_app =()=>{
        setApp(app);
    }

    return(
        <li key={app.id} onClick={set_as_currrent_app}>
            <div>{app.name}</div>
            <div>Rate Limit: {app.rate_limit}</div>
            <div>Refill Rate: {app.refill_rate}</div>
        </li>
    )
}