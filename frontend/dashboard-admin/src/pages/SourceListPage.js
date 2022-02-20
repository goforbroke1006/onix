import React from "react";

class SourceListPage extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            sources: [],
        }
    }

    componentDidMount() {
        this.loadSourcesList()
    }

    render() {
        return (
            <div>
                <h1>Source list Page</h1>
                <div>
                    <table>
                        <thead>
                        <tr>
                            <th>ID</th>
                            <th>Name</th>
                            <th>Kind</th>
                            <th>Address</th>
                        </tr>
                        </thead>

                        <tbody>
                        {this.state.sources.map((item, index) => {
                            return (
                                <tr key={"source-list-item-" + index}>
                                    <td>{item.id}</td>
                                    <td>{item.title}</td>
                                    <td>{item.kind}</td>
                                    <td>{item.address}</td>
                                    <td>
                                        <button>edit</button>
                                        <button>remove</button>
                                    </td>
                                </tr>
                            )
                        })}
                        </tbody>
                    </table>
                </div>
            </div>
        )
    }

    loadSourcesList = () => {
        this.props.provider.loadSources()
            .then(sources => this.setState({sources: sources}));
    }
}

export default SourceListPage;
