import React from "react";

import DashboardMainApiClient from "./external/client";

import {Button, Col, Row} from "antd";
import ServiceDropDown from "./component/ServiceDropDown";
import SourceDropDown from "./component/SourceDropDown";
import ReleaseDropDown from "./component/ReleaseDropDown";
import DateTimePickerWithLimits from "./component/DateTimePickerWithLimits";
import PeriodDropDown from "./component/PeriodDropDown";
import CompareReleasesPanel from "./component/CompareReleasesPanel";

import 'antd/dist/antd.dark.css';
import './App.css';
import CompareReleasesLinks from "./component/CompareReleasesLinks";

class App extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
            serviceTitle: null,

            releaseOne: null,
            releaseOneSourceId: null,
            releaseOneStartAt: null,

            releaseTwo: null,
            releaseTwoSourceId: null,
            releaseTwoStartAt: null,

            period: null,

            comparison: [],
        };

        this.client = new DashboardMainApiClient();
    }

    render() {
        return (
            <div>

                <Row gutter={[16, 16]}>
                    <Col span={6}/>
                    <Col span={12}>

                        <Row gutter={[16, 16]}>
                            <Col span={24}><h3>Service</h3></Col>
                            <Col span={6}><label>Service name:</label></Col>
                            <Col span={18}><ServiceDropDown provider={this.client}
                                                            onChange={this.onServiceSelected}/></Col>
                        </Row>

                        <Row gutter={[16, 16]}>
                            <Col span={24}><h3>Release 1</h3></Col>

                            <Col span={6}><label>Source:</label></Col>
                            <Col span={18}><SourceDropDown provider={this.client}
                                                           onChange={this.onReleaseOneSourceSelected}/></Col>

                            <Col span={6}><label>Release name:</label></Col>
                            <Col span={18}><ReleaseDropDown provider={this.client}
                                                            serviceName={this.state.serviceTitle}
                                                            onChange={this.onReleaseOneSelected}/></Col>

                            <Col span={6}><label>Start:</label></Col>
                            <Col span={18}><DateTimePickerWithLimits
                                from={this.state.releaseOne ? this.state.releaseOne.from : 0}
                                till={this.state.releaseOne ? this.state.releaseOne.till : 0}
                                onChange={this.onReleaseOneSelectTimeRange}/></Col>
                        </Row>

                        <Row gutter={[16, 16]}>
                            <Col span={24}><h3>Release 2</h3></Col>

                            <Col span={6}><label>Source:</label></Col>
                            <Col span={18}><SourceDropDown provider={this.client}
                                                           onChange={this.onReleaseTwoSourceSelected}/></Col>

                            <Col span={6}><label>Release name:</label></Col>
                            <Col span={18}><ReleaseDropDown provider={this.client}
                                                            serviceName={this.state.serviceTitle}
                                                            onChange={this.onReleaseTwoSelected}/></Col>

                            <Col span={6}><label>Start:</label></Col>
                            <Col span={18}><DateTimePickerWithLimits
                                from={this.state.releaseTwo ? this.state.releaseTwo.from : 0}
                                till={this.state.releaseTwo ? this.state.releaseTwo.till : 0}
                                onChange={this.onReleaseTwoSelectTimeRange}/></Col>
                        </Row>

                        <Row gutter={[16, 16]}>
                            <Col span={6}><label>Period:</label></Col>
                            <Col span={18}><PeriodDropDown onChange={this.onPeriodSelected}/></Col>
                        </Row>
                    </Col>
                    <Col span={6}/>
                </Row>


                <Row gutter={[16, 16]}>
                    <Col span={10}/>
                    <Col span={4}>
                        <Button onClick={this.onLoadMetricsClick}>Load</Button>
                    </Col>
                    <Col span={10}/>
                </Row>


                <Row gutter={[16, 16]}>
                    <Col span={2}/>
                    <Col span={4}>
                        <CompareReleasesLinks comparison={this.state.comparison}/>
                    </Col>
                    <Col span={16}>
                        <CompareReleasesPanel comparison={this.state.comparison}/>
                    </Col>
                    <Col span={2}/>
                </Row>


            </div>
        );
    }

    /**
     *
     * @param {String} serviceTitle
     */
    onServiceSelected = (serviceTitle) => {
        this.setState({serviceTitle: serviceTitle});
    };

    /**
     *
     * @param {number} sourceId
     */
    onReleaseOneSourceSelected = (sourceId) => {
        this.setState({releaseOneSourceId: sourceId});
    }

    /**
     *
     * @param {number} sourceId
     */
    onReleaseTwoSourceSelected = (sourceId) => {
        this.setState({releaseTwoSourceId: sourceId});
    }

    onReleaseOneSelected = (release) => {
        console.log("set release one");
        console.log(release);
        this.setState({releaseOne: release});
    };

    onReleaseTwoSelected = (release) => {
        console.log("set release two");
        console.log(release);
        this.setState({releaseTwo: release});
    };

    /**
     *
     * @param {String} period
     */
    onPeriodSelected = (period) => {
        this.setState({period: period});
    }

    /**
     *
     * @param {number} unix
     */
    onReleaseOneSelectTimeRange = (unix) => {
        this.setState({releaseOneStartAt: unix});
    }

    /**
     *
     * @param {number} unix
     */
    onReleaseTwoSelectTimeRange = (unix) => {
        this.setState({releaseTwoStartAt: unix});
    }

    onLoadMetricsClick = () => {
        if (!this.state.serviceTitle) return;

        if (!this.state.releaseOne) return;
        if (!this.state.releaseOneStartAt) return;
        if (!this.state.releaseOneSourceId) return;

        if (!this.state.releaseTwo) return;
        if (!this.state.releaseTwoStartAt) return;
        if (!this.state.releaseTwoSourceId) return;

        if (!this.state.period) return;

        this.client.loadComparison(
            this.state.serviceTitle,
            this.state.releaseOne.title,
            this.state.releaseOneStartAt,
            this.state.releaseOneSourceId,
            this.state.releaseTwo.title,
            this.state.releaseTwoStartAt,
            this.state.releaseTwoSourceId,
            this.state.period
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
                    switch (this.state.period) {
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

            this.setState({comparison: stateReports})
        })
    }
}

export default App;
