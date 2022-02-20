import React from "react";

class MainPage extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            services: [],
        };
    }

    componentDidMount() {
        this.loadServiceList()
    }

    render() {
        return (
            <div>
                <h1>Main Page</h1>
                <div>
                    <h2>Services under investigation</h2>

                    <ul>
                        {this.state.services.map((service, index) => {
                            return (
                                <li key={"services-list-item-" + index}>
                                    {service.title}
                                    &nbsp;
                                    &nbsp;
                                    &nbsp;
                                    ({service.releases.map((releaseTitle, rti) => {
                                        return (
                                            <span><span>{releaseTitle}</span>, &nbsp;</span>
                                        )
                                    })})
                                </li>
                            )
                        })}
                    </ul>
                </div>
            </div>
        )
    }

    loadServiceList = () => {
        this.props.provider.loadServices()
            .then(services => this.setState({services: services}));
    }
}

export default MainPage;
