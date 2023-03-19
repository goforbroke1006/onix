import React from "react";
import PropTypes from 'prop-types';
import {Chart} from "react-google-charts";
import {getAverageDiff} from "../../math/series";

class CompareReleasesChart extends React.Component {
    static propTypes = {
        criteriaName: PropTypes.any,
        releaseOne: PropTypes.string,
        releaseTwo: PropTypes.string,
        graph: PropTypes.arrayOf(PropTypes.shape({
            t1: PropTypes.string,
            t2: PropTypes.string,
            v1: PropTypes.number,
            v2: PropTypes.number,
        })),
        direction: PropTypes.string,
    };

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

        const diff = getAverageDiff(values1, values2);

        let color = "#00ff00";
        switch (this.props.direction) {
            case "increase":
                if (diff.absolute < 0) color = "#ff0000";
                break
            case "decrease":
                if (diff.absolute > 0) color = "#ff0000";
                break
            case "equal":
                color = "#999999"
        }

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
                <div style={{color: color}}>Absolute: {diff.absolute}</div>
                <div style={{color: color}}>Percentage: {diff.percentage}</div>
            </div>
        );
    }
}


export default CompareReleasesChart;
