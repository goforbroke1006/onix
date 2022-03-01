export default class DashboardMainApiClient {
    constructor() {
        this.baseUrl = process.env.REACT_APP_API_DASHBOARD_MAIN_BASE_ADDR ?? 'http://127.0.0.1:8082/api/dashboard-main';
    }

    loadServices() {
        return fetch(`${this.baseUrl}/service`)
            .then(response => response.json());
    }

    loadSourcesList() {
        return fetch(`${this.baseUrl}/source`)
            .then(response => response.json());
    }

    loadReleasesList(serviceName) {
        return fetch(`${this.baseUrl}/release?service=${serviceName}`)
            .then((response) => {
                if (response.ok)
                    return response.json();
                else
                    return [];
            })
    }

    /**
     *
     * @param serviceTitle
     * @param releaseOneTitle
     * @param releaseOneStartAt
     * @param releaseOneSourceId
     * @param releaseTwoTitle
     * @param releaseTwoStartAt
     * @param releaseTwoSourceId
     * @param period
     * @returns {Promise<Response>}
     */
    loadComparison(
        serviceTitle,
        releaseOneTitle, releaseOneStartAt, releaseOneSourceId,
        releaseTwoTitle, releaseTwoStartAt, releaseTwoSourceId,
        period
    ) {
        let url = `${this.baseUrl}/compare?service=${serviceTitle}` +
            `&release_one_title=${releaseOneTitle}` +
            `&release_one_start=${releaseOneStartAt}` +
            `&release_one_source_id=${releaseOneSourceId}` +
            `&release_two_title=${releaseTwoTitle}` +
            `&release_two_start=${releaseTwoStartAt}` +
            `&release_two_source_id=${releaseTwoSourceId}` +
            `&period=${period}`
        return fetch(url)
            .then(response => {
                if (response.ok)
                    return response.json()
                else {
                    console.warn('load comparison failed')
                    return [];
                }
            });
    }
}