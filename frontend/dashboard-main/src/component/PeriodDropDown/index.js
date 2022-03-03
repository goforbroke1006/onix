import React from "react";
import PropTypes from "prop-types";

const defaultValue = "1h";

export default class PeriodDropDown extends React.Component {
    static propTypes = {
        onChange: PropTypes.func,
        provider: PropTypes.object,
    };

    constructor(props) {
        super(props);

        this.state = {
            items: ["15m", "1h", "6h", "1d"],
            value: defaultValue,
        }

        this.props.onChange(defaultValue);
    }

    render() {
        return (
            <select onChange={this.onChange} value={this.state.value}>
                {this.state.items.map((period) => {
                    return (
                        <option key={"period-" + period} value={"" + period}>{period}</option>
                    )
                })}
            </select>
        )
    }

    onChange = (event) => {
        let selected = event.target.value;

        this.setState({value: selected});

        if (this.props.onChange) {
            this.props.onChange(selected);
        }
    }
}
