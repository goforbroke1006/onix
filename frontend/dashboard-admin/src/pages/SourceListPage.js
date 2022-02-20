import React from "react";

class SourceListPage extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            items: [],
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
                        {this.state.items.map((item, index) => {
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
        let baseUrl = process.env.REACT_APP_API_DASHBOARD_ADMIN_BASE_ADDR ?? 'http://127.0.0.1:8083/api/dashboard-admin';
        fetch(baseUrl + "/source")
            .then(response => response.json())
            .then(data => this.setState({items: data}))
    }
}

export default SourceListPage;
