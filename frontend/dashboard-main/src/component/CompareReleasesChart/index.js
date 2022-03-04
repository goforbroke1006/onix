import React from "react";
import PropTypes from 'prop-types';
import {Chart} from "react-google-charts";

class CompareReleasesChart extends React.Component {
    static propTypes = {
        title: PropTypes.any,
        measurements: PropTypes.array,
    };

    constructor(props) {
        super(props);
        this.state = {
            title: props.title,
            measurements: props.measurements,
        }
    }

    componentDidUpdate(prevProps) {
        if (prevProps.measurements !== this.props.measurements) {
            this.setState({
                measurements: this.props.measurements,
            })
        }
    }

    render() {
        if (!Array.isArray(this.props.measurements)) {
            return (
                <div>--- no data ---</div>
            )
        }

        // const dateOne = new Date(this.state.releaseOneStartAt * 1000);
        // const dateTwo = new Date(this.state.releaseTwoStartAt * 1000);

        let options = {
//             title: `Releases performance by "${this.props.title}" criteria \n
// Release 1: ${dateOne.toUTCString()}
// Release 2: ${dateTwo.toUTCString()}`,
            curveType: 'function',
            legend: {position: 'bottom'},
            hAxis: {
                direction: 1,
                slantedText: true,
                slantedTextAngle: 90,
            }
        };

        return (
            <div>
                <Chart
                    chartType="LineChart"
                    data={this.props.measurements}
                    width="100%"
                    height="500px"
                    legendToggle
                    options={options}
                    formatters={[
                        {
                            type: "NumberFormat",
                            column: 1,
                            options: {pattern: '##.#######'}
                        },
                        {
                            type: "NumberFormat",
                            column: 2,
                            options: {pattern: '##.#######'}
                        }
                    ]}
                />
            </div>
        );
    }
}

export default CompareReleasesChart;
