import React from "react";
import CompareReleasesChart from "../CompareReleasesChart";
import PropTypes from "prop-types";

class CompareReleasesPanel extends React.Component {
    static propTypes = {
        provider: PropTypes.object,

        serviceTitle: PropTypes.string,

        releaseOneTitle: PropTypes.string,
        releaseOneStartAt: PropTypes.any,
        releaseOneSourceId: PropTypes.number,

        releaseTwoTitle: PropTypes.string,
        releaseTwoStartAt: PropTypes.any,
        releaseTwoSourceId: PropTypes.number,

        period: PropTypes.string,
    };

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

        if (prevProps.releaseOneStartAt !== this.props.releaseOneStartAt
            || prevProps.releaseOneSourceId !== this.props.releaseOneSourceId
            || prevProps.releaseTwoStartAt !== this.props.releaseTwoStartAt
            || prevProps.releaseTwoSourceId !== this.props.releaseTwoSourceId
            || prevProps.period !== this.props.period) {
            this.loadCompareReport();
        }
    }

    render() {
        return (
            <div>
                <ul>
                    {this.state.reports.map((report, index) => {
                        return (
                            <li key={`jump-to-${index}`}><a href={"#" + report.title}>{report.title}</a></li>
                        )
                    })}
                </ul>
                <div>
                    {this.state.reports.map((report, index) => {
                        return (
                            <div key={`wrap-chart-${index}`}>
                                <span id={report.title}>&nbsp;</span>
                                <CompareReleasesChart
                                    key={`chart-${index}`}
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

        if (!this.props.releaseOneTitle) return;
        if (!this.props.releaseOneStartAt) return;
        if (!this.props.releaseOneSourceId) return;

        if (!this.props.releaseTwoTitle) return;
        if (!this.props.releaseTwoStartAt) return;
        if (!this.props.releaseTwoSourceId) return;

        if (!this.props.period) return;

        this.props.provider.loadComparison(
            this.props.serviceTitle,
            this.props.releaseOneTitle,
            this.props.releaseOneStartAt,
            this.props.releaseOneSourceId,
            this.props.releaseTwoTitle,
            this.props.releaseTwoStartAt,
            this.props.releaseTwoSourceId,
            this.props.period
        ).then(resp => {

            let stateReports = [];

            for (let ri = 0; ri < resp.reports.length; ri++) {
                let report = resp.reports[ri];

                let measurements = [
                    ['Time', resp.release_one, resp.release_two],
                ];
                for (let gi = 0; gi < report.graph.length; gi++) {
                    let graphEl = report.graph[gi];

                    let timeLabel = `${graphEl.t1}\n${graphEl.t2}`;
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