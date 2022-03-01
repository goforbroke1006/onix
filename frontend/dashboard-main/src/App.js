import React from "react";

import './App.css';

import ServiceDropDown from "./component/ServiceDropDown";
import ReleaseDropDown from "./component/ReleaseDropDown";
import CompareReleasesPanel from "./component/CompareReleasesPanel";
import DateTimePickerWithLimits from "./component/DateTimePickerWithLimits";
import PeriodDropDown from "./component/PeriodDropDown";
import SourceDropDown from "./component/SourceDropDown";
import DashboardMainApiClient from "./external/client";

class App extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
            serviceTitle: null,

            releaseOneTitle: null,
            releaseOneFrom: null,
            releaseOneTill: null,
            releaseOneSourceId: null,

            releaseTwoTitle: null,
            releaseTwoFrom: null,
            releaseTwoTill: null,
            releaseTwoSourceId: null,
        };
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
        // console.log("set release one");
        // console.log(release);

        this.setState({
            releaseOneTitle: release.title,
            releaseOneFrom: release.from,
            releaseOneTill: release.till,
            releaseOneStartAt: release.from,
        });
    };

    onReleaseTwoSelected = (release) => {
        // console.log("set release two");
        // console.log(release);

        this.setState({
            releaseTwoTitle: release.title,
            releaseTwoFrom: release.from,
            releaseTwoTill: release.till,
            releaseTwoStartAt: release.from,
        });
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

    render() {
        const client = new DashboardMainApiClient();

        return (
            <div className="App">
                <ServiceDropDown provider={client} onChange={this.onServiceSelected}/>

                <ReleaseDropDown provider={client}
                                 serviceName={this.state.serviceTitle}
                                 onChange={this.onReleaseOneSelected}/>
                <SourceDropDown provider={client} onChange={this.onReleaseOneSourceSelected}/>
                <DateTimePickerWithLimits from={this.state.releaseOneFrom} till={this.state.releaseOneTill}
                                          onChange={this.onReleaseOneSelectTimeRange}/>

                <ReleaseDropDown provider={client}
                                 serviceName={this.state.serviceTitle}
                                 onChange={this.onReleaseTwoSelected}/>
                <SourceDropDown provider={client} onChange={this.onReleaseTwoSourceSelected}/>
                <DateTimePickerWithLimits from={this.state.releaseTwoFrom} till={this.state.releaseTwoTill}
                                          onChange={this.onReleaseTwoSelectTimeRange}/>

                <PeriodDropDown onChange={this.onPeriodSelected}/>

                <CompareReleasesPanel provider={client}

                                      serviceTitle={this.state.serviceTitle}

                                      releaseOneTitle={this.state.releaseOneTitle}
                                      releaseOneStartAt={this.state.releaseOneStartAt}
                                      releaseOneSourceId={this.state.releaseOneSourceId}

                                      releaseTwoTitle={this.state.releaseTwoTitle}
                                      releaseTwoStartAt={this.state.releaseTwoStartAt}
                                      releaseTwoSourceId={this.state.releaseTwoSourceId}

                                      period={this.state.period}
                />

            </div>
        );
    }
}

export default App;
