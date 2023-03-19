import React from "react";
import CompareReleasesChart from "../CompareReleasesChart";
import PropTypes from "prop-types";

class CompareReleasesPanel extends React.Component {
    static propTypes = {
        comparison: PropTypes.array,
    };

    render() {
        if (!Array.isArray(this.props.comparison)) {
            return (
                <div>--- no data ---</div>
            )
        }

        return (
            <div>

                <div>
                    {this.props.comparison.map((report, index) => {
                        return (
                            <div key={`wrap-chart-${index}`}>
                                <span id={report.title}>&nbsp;</span>
                                <CompareReleasesChart
                                    key={`chart-${index}`}
                                    criteriaName={report.criteriaName}
                                    releaseOne={report.releaseOne}
                                    releaseTwo={report.releaseTwo}
                                    graph={report.graph}
                                    direction={report.direction}
                                />
                            </div>
                        )
                    })}
                </div>
            </div>
        )
    }
}

export default CompareReleasesPanel;