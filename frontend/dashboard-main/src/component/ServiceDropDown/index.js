import React from "react";

class ServiceDropDown extends React.Component {
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
                                    value={object.title}>
                                {object.title}
                            </option>
                        )
                    })}
                </select>
            </div>
        )
    }

     loadServices = () => {
        let baseUrl = process.env.REACT_APP_API_DASHBOARD_MAIN_BASE_ADDR ?? 'http://127.0.0.1:8082/api/dashboard-main';
        fetch(baseUrl + "/service")
            .then(response => response.json())
            .then(data => this.setState({items: data}))
    }
}

export default ServiceDropDown;