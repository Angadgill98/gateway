


export class Api {
    url:string
    constructor(url:string){
        this.url=url
    }

    async RegisterApp(app_name: string): Promise<ApiResponse<any>> {
        const route = "/register-app";

        const body = {
            name: app_name
        };

        console.log("Request:", route, body);

        try {
            const res = await fetch(this.url + route, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify(body)
            });

            const data = await res.json();

            console.log("Response:", res.status, data);

            return {
                success: res.ok && data.success,
                message: data.message,
                data: data.data ?? null
            };
        } catch (error) {
            console.error("RegisterApp failed:", error);

            return {
                success: false,
                message: error instanceof Error ? error.message : "Request failed",
                data: null
            };
        }
    }

    async GetAllApps(): Promise<ApiResponse<any[]>> {
        const route = "/get-all-apps";

        try {
            const res = await fetch(this.url + route, {
                method: "GET",
                headers: {
                    "Content-Type": "application/json"
                }
            });

            const data = await res.json();

            console.log("Response:", res.status, data);

            return {
                success: res.ok && data.success,
                message: data.message,
                data: data.data ?? null
            };
        } catch (error) {
            console.error("GetAllApps failed:", error);

            return {
                success: false,
                message: error instanceof Error ? error.message : "Request failed",
                data: null
            };
        }
    }
}

interface ApiResponse<T = null> {
    success: boolean;
    message: string;
    data: T | null;
}


const api = new Api("http://localhost:8080/api");

export default api;