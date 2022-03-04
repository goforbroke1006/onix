import React from "react";
import PropTypes from "prop-types";
import {Select} from "antd";

const {Option} = Select;

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
        let emptyOption;
        if (this.state.items.length === 0) {
            emptyOption = (<Option value={""}>-- no data --</Option>)
        } else {
            emptyOption = (<Option value={""}>-- none --</Option>)
        }

        return (
            <Select defaultValue={this.state.value} style={{width: 600}} onChange={this.onChange}>
                {emptyOption}
                {this.state.items.map((period) => {
                    return (
                        <Option key={`period-${period}`} value={"" + period}>{period}</Option>
                    )
                })}
            </Select>
        )
    }

    onChange = (period) => {
        this.setState({value: period});

        if (this.props.onChange) {
            this.props.onChange(period);
        }
    }
}
