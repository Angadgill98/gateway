export interface App {
    id: string;
    name: string;
    rate_limit: number;
    refill_rate: number;
}

export interface Endpoint {
    id: string;
    app_id: string;
    route: string;
    destination_ids: string[];
    rate_limit: number;
    refill_rate: number;
}