import React from "react";
import PropTypes from 'prop-types';
import {Chart} from "react-google-charts";
import {getAverageDiff} from "../../math/series";

class CompareReleasesChart extends React.Component {
    static propTypes = {
        criteriaName: PropTypes.any,
        // measurements: PropTypes.array,
        releaseOne: PropTypes.string,
        releaseTwo: PropTypes.string,
        graph: PropTypes.array,
    };

    // constructor(props) {
    //     super(props);
    //     this.state = {
    //         title: props.title,
    //         // measurements: props.measurements,
    //         releaseOne: props.releaseOne,
    //         releaseTwo: props.releaseTwo,
    //         graph: props.graph,
    //     }
    // }

    // componentDidUpdate(prevProps, prevState, snapshot) {
    //     if (prevProps.seriesOne !== this.props.seriesOne) {
    //         this.setState({seriesOne: this.props.seriesOne,})
    //     }
    //     if (prevProps.seriesTwo !== this.props.seriesTwo) {
    //         this.setState({seriesTwo: this.props.seriesTwo,})
    //     }
    // }

    render() {
        if (!Array.isArray(this.props.graph)) {
            return (
                <div>--- no data ---</div>
            )
        }

        let chartData = [
            ['Time', this.props.releaseOne, this.props.releaseTwo],
        ];
        let values1 = [];
        let values2 = [];

        for (let gi = 0; gi < this.props.graph.length; gi++) {
            let graphEl = this.props.graph[gi];

            let timeLabel = `${graphEl.t1}\n${graphEl.t2}`;
            let v1 = graphEl.v1;
            let v2 = graphEl.v2;
            chartData.push([timeLabel, v1, v2]);

            values1.push(v1);
            values2.push(v2);
        }

        // load zero (null-as-zero logic)
        if (this.props.graph.length === 0) {
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
                chartData.push(["" + fmi, 0, 0]);
            }
        }

        const diff = getAverageDiff(values1, values2);

        let options = {
            title: `Releases performance by "${this.props.criteriaName}" criteria`,
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
                    data={chartData}
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
                <div>Absolute: {diff.absolute}</div>
                <div>Percentage: {diff.percentage}</div>
            </div>
        );
    }
}

export default CompareReleasesChart;
