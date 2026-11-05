import { createContext, useContext, useState, type ReactNode } from "react";
import type { App, Endpoint } from "../models";



interface GlobalContextType {
    app: App | null;
    setApp: React.Dispatch<React.SetStateAction<App | null>>;

    endpoints: Endpoint[];
    setEndpoints: React.Dispatch<React.SetStateAction<Endpoint[]>>;
}


const Global_context = createContext<GlobalContextType | undefined>(undefined);

export function GlobalProvider({ children }: { children: ReactNode }) {
    const [app, setApp] = useState<App | null>(null);
    const [endpoints, setEndpoints] = useState<Endpoint[]>([]);

    return (
        <Global_context.Provider value={{ app, setApp, endpoints, setEndpoints }}>
            {children}
        </Global_context.Provider>
    );
}

export function useGlobal() {
    const context = useContext(Global_context);

    if (!context) {
        throw new Error("useGlobal must be used within GlobalProvider");
    }

    return context;
}