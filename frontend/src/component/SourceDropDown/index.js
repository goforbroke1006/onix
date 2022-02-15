import React from "react";

class SourceDropDown extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            items: [],
            value: "",
        };
    }

    componentDidMount() {
        this.loadSources();
    }

    onChange = (event) => {
        let sourceId = event.target.value;

        this.setState({value: sourceId});

        if (this.props.onChange)
            this.props.onChange(sourceId);
    }

    render() {
        return (
            <select onChange={this.onChange} value={this.state.value}>
                <option>----</option>
                {this.state.items.map((object, index) => {
                    return (
                        <option key={"source-" + index} value={object.id}>
                            {object.title} [{object.kind}] {object.address}
                        </option>
                    )
                })}
            </select>
        )
    }

    loadSources = () => {
        let baseUrl = process.env.REACT_APP_API_DASHBOARD_MAIN_BASE_ADDR ?? 'http://127.0.0.1:8082/api/dashboard-main';
        fetch(baseUrl + "/source")
            .then(response => response.json())
            .then(data => this.setState({items: data}))
    }
}

export default SourceDropDown;