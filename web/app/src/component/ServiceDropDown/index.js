import React from "react";
import PropTypes from "prop-types";
import {Select} from "antd";

const {Option} = Select;

export default class ServiceDropDown extends React.Component {
    static propTypes = {
        onChange: PropTypes.func,
        provider: PropTypes.object,
    };

    constructor(props) {
        super(props);
        this.state = {
            items: [],
            value: "",
        }
    }

    componentDidMount() {
        this.loadServices();
    }

    onChange = (serviceName) => {
        console.log(serviceName);

        this.setState({value: serviceName});

        if (this.props.onChange)
            this.props.onChange(serviceName);
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
                {this.state.items.map((object, index) => {
                    return (
                        <Option key={"service-select-option-" + index}
                                value={object.title}>{object.title}</Option>
                    )
                })}
            </Select>
        )
    }

    loadServices = () => {
        this.props.provider.loadServices()
            .then(data => this.setState({items: data}))
    }
}
