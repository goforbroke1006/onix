import React from "react";

import './App.css';

import ServiceDropDown from "./component/ServiceDropDown";
import ReleaseDropDown from "./component/ReleaseDropDown";
import CompareReleasesPanel from "./component/CompareReleasesPanel";
import DateTimePickerWithLimits from "./component/DateTimePickerWithLimits";
import PeriodDropDown from "./component/PeriodDropDown";
import SourceDropDown from "./component/SourceDropDown";

class App extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
            serviceTitle: null,

            releaseOneTitle: null,
            releaseOneFrom: null,
            releaseOneTill: null,

            releaseTwoTitle: null,
            releaseTwoFrom: null,
            releaseTwoTill: null,
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
    onSourceSelected = (sourceId) => {
        this.setState({sourceId: sourceId});
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
        return (
            <div className="App">
                <ServiceDropDown onChange={this.onServiceSelected}/>

                <SourceDropDown onChange={this.onSourceSelected}/>

                <ReleaseDropDown serviceName={this.state.serviceTitle}
                                 sourceId={this.state.sourceId}
                                 onChange={this.onReleaseOneSelected}/>
                <ReleaseDropDown serviceName={this.state.serviceTitle}
                                 sourceId={this.state.sourceId}
                                 onChange={this.onReleaseTwoSelected}/>

                <DateTimePickerWithLimits from={this.state.releaseOneFrom} till={this.state.releaseOneTill}
                                          onChange={this.onReleaseOneSelectTimeRange}/>
                <DateTimePickerWithLimits from={this.state.releaseTwoFrom} till={this.state.releaseTwoTill}
                                          onChange={this.onReleaseTwoSelectTimeRange}/>

                <PeriodDropDown onChange={this.onPeriodSelected}/>

                <CompareReleasesPanel
                    serviceTitle={this.state.serviceTitle}
                    sourceId={this.state.sourceId}
                    releaseOneTitle={this.state.releaseOneTitle}
                    releaseOneStartAt={this.state.releaseOneStartAt}
                    releaseTwoTitle={this.state.releaseTwoTitle}
                    releaseTwoStartAt={this.state.releaseTwoStartAt}
                    period={this.state.period}
                />

            </div>
        );
    }
}

export default App;
