import React from "react";
import PropTypes from "prop-types";

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

    onChange = (event) => {
        let serviceName = event.target.value;

        console.log(serviceName);

        this.setState({value: serviceName});

        if (this.props.onChange)
            this.props.onChange(serviceName);
    }

    render() {
        return (
            <div>
                <select onChange={this.onChange} value={this.state.value}>
                    <option>----</option>
                    {this.state.items.map((object, index) => {
                        return (
                            <option key={"service-select-option-" + index}
                                    value={object.title}>{object.title}</option>
                        )
                    })}
                </select>
            </div>
        )
    }

    loadServices = () => {
        this.props.provider.loadServices().then(data => this.setState({items: data}))
    }
}
