class ApiDataProvider {
    constructor() {
        const localhostApiAddr = 'http://127.0.0.1:8083/api/dashboard-admin';
        this.baseAddr = process.env.REACT_APP_API_DASHBOARD_ADMIN_BASE_ADDR ?? localhostApiAddr;
    }

    loadServices = async () => {
        return await fetch(`${this.baseAddr}/service`)
            .then(response => {
                if (!response.ok) {
                    return [];
                }
                return response.json()
            })
            .then(data => {
                return data;
            })
            .catch((err) => console.warn(err));
    }

    loadSources = async () => {
        return await fetch(`${this.baseAddr}/source`)
            .then(response => {
                if (!response.ok) {
                    return [];
                }
                return response.json();
            })
            .then(data => {
                return data;
            })
            .catch((err) => console.warn(err));
    }
}

export default ApiDataProvider;