import React from "react";
import CompareReleasesChart from "../CompareReleasesChart";

class CompareReleasesPanel extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
            reports: [],
        }
    }

    componentDidMount() {
        this.loadCompareReport();
    }

    componentDidUpdate(prevProps) {
        if (prevProps.sourceId !== this.props.sourceId) {
            this.loadCompareReport();
        }
        if (prevProps.releaseOneStartAt !== this.props.releaseOneStartAt) {
            this.loadCompareReport();
        }
        if (prevProps.releaseTwoStartAt !== this.props.releaseTwoStartAt) {
            this.loadCompareReport();
        }
        if (prevProps.period !== this.props.period) {
            this.loadCompareReport();
        }
    }

    render() {
        return (
            <div>
                <ul>
                    {this.state.reports.map((report, index) => {
                        return (
                            <li><a href={"#" + report.title}>{report.title}</a></li>
                        )
                    })}
                </ul>
                <div>
                    {this.state.reports.map((report, index) => {
                        return (
                            <div>
                                <span id={report.title}>&nbsp;</span>
                                <CompareReleasesChart key={"chart-" + index}
                                                      title={report.title}
                                                      measurements={report.measurements}
                                                      releaseOneStartAt={this.props.releaseOneStartAt}
                                                      releaseTwoStartAt={this.props.releaseTwoStartAt}
                                />
                            </div>
                        )
                    })}
                </div>
            </div>
        )
    }

    loadCompareReport = () => {
        if (!this.props.serviceTitle) return;
        if (!this.props.sourceId) return;
        if (!this.props.releaseOneTitle) return;
        if (!this.props.releaseTwoTitle) return;
        if (!this.props.releaseOneStartAt) return;
        if (!this.props.releaseTwoStartAt) return;
        if (!this.props.period) return;

        let baseUrl = process.env.API_DASHBOARD_MAIN_BASE_ADDR ?? 'http://127.0.0.1:8082/api/dashboard-main';

        let url = baseUrl + `/compare?service=${this.props.serviceTitle}&source_id=${this.props.sourceId}` +
            `&release_one_title=${this.props.releaseOneTitle}` +
            `&release_one_start=${this.props.releaseOneStartAt}` +
            `&release_two_title=${this.props.releaseTwoTitle}` +
            `&release_two_start=${this.props.releaseTwoStartAt}` +
            `&period=${this.props.period}`
        fetch(url)
            .then(response => response.json())
            .then(extractResp => {

                let stateReports = [];

                for (let ri = 0; ri < extractResp.reports.length; ri++) {
                    let report = extractResp.reports[ri];

                    let measurements = [
                        ['Time', extractResp.release_one, extractResp.release_two],
                    ];
                    for (let gi = 0; gi < report.graph.length; gi++) {
                        let graphEl = report.graph[gi];

                        let timeLabel = `# ${gi} - ${graphEl.t1} - ${graphEl.t2}`;
                        let v1 = graphEl.v1;
                        let v2 = graphEl.v2;
                        measurements.push([timeLabel, v1, v2]);
                    }

                    // load zero (null-as-zero logic)
                    if (report.graph.length === 0) {
                        let countFakeItems = 0;
                        switch (this.props.period) {
                            case '15m':
                                countFakeItems = 15 / 5;
                                break;
                            case '1h':
                                countFakeItems = 60 / 5;
                                break;
                            case '6h':
                                countFakeItems = 6 * 60 / 5;
                                break;
                            case '1d':
                                countFakeItems = 24 * 60 / 5;
                                break;
                        }
                        for (let fmi = 0; fmi < countFakeItems; fmi++) {
                            measurements.push(["" + fmi, 0, 0]);
                        }
                    }

                    stateReports.push({
                        title: report.title,
                        measurements: measurements,
                    });
                }

                this.setState({reports: stateReports})
            })
    }
}

export default CompareReleasesPanel;