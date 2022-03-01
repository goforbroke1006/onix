import React from "react";

export default class SourceDropDown extends React.Component {
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
        let sourceId = parseInt(event.target.value);

        this.setState({value: sourceId});

        if (this.props.onChange)
            this.props.onChange(sourceId);
    }

    render() {
        if (this.state.items.length === 0) {
            return (
                <select onChange={this.onChange} value={this.state.value}>
                    <option>-- no data --</option>
                </select>
            )
        }

        return (
            <select onChange={this.onChange} value={this.state.value}>
                <option>-- none --</option>
                {this.state.items.map((object, index) => {
                    return (
                        <option key={`source-${index}`}
                                value={object.id}>{`${object.title} [${object.kind}] ${object.address}`}</option>
                    )
                })}
            </select>
        )
    }

    loadSources = () => {
        return this.props.provider.loadSourcesList()
            .then(data => {
                this.setState({items: data});
            })
            .catch((err) => {
                console.warn(err);
                this.setState({items: []});
            })
    }
}
