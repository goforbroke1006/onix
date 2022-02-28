export default class DashboardMainApiClient {
    constructor() {
        this.baseUrl = process.env.REACT_APP_API_DASHBOARD_MAIN_BASE_ADDR ?? 'http://127.0.0.1:8082/api/dashboard-main';
    }

    async loadServices() {
        return await fetch(`${this.baseUrl}/service`)
            .then(response => response.json());
    }

    async loadSourcesList() {
        return await fetch(`${this.baseUrl}/source`)
            .then(response => response.json());
    }

    async loadReleasesList(serviceName) {
        return await fetch(`${this.baseUrl}/release?service=${serviceName}`)
            .then(response => response.json())
    }
}