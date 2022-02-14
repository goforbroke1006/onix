import React from "react";

class ReleaseDropDown extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            items: [],
        };
    }

    componentDidMount() {
        this.loadRelease();
    }

    componentDidUpdate(prevProps) {
        if (prevProps.serviceName !== this.props.serviceName) {
            this.loadRelease();
        }
        if (prevProps.sourceId !== this.props.sourceId) {
            this.loadRelease();
        }
    }

    render() {
        return (
            <div>
                <select onChange={this.onChange} value={this.state.value}>
                    {this.state.items.map(function (object, index) {
                        return (
                            <option key={"service-select-option-" + object.id}
                                    value={object.id}>
                                {object.title}
                            </option>
                        )
                    })}
                </select>
            </div>
        )
    }

    onChange = (event) => {
        let releaseId = event.target.value;

        this.setState({value: releaseId});

        let release = null;
        for (let i = 0; i < this.state.items.length; i++) {
            if (parseInt(this.state.items[i].id) === parseInt(releaseId)) {
                release = this.state.items[i];
                break;
            }
        }

        if (null != release && this.props.onChange)
            this.props.onChange(release);
    }

    loadRelease = () => {
        if (!this.props.serviceName) return;
        if (!this.props.sourceId) return;

        let baseUrl = process.env.API_DASHBOARD_MAIN_BASE_ADDR ?? 'http://127.0.0.1:8082/api/dashboard-main';
        fetch(`${baseUrl}/release-summary?service=${this.props.serviceName}&source_id=${this.props.sourceId}`)
            .then(response => response.json())
            .then(data => this.setState({items: data}))
    }
}

export default ReleaseDropDown;