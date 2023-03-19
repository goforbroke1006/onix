import React from "react";
import PropTypes from "prop-types";

export default class CompareReleasesLinks extends React.Component {
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
            <ul>
                {this.props.comparison.map((report, index) => {
                    return (
                        <li key={`jump-to-${index}`}><a href={"#" + report.criteriaName}>{report.criteriaName}</a></li>
                    )
                })}
            </ul>
        );
    }
}