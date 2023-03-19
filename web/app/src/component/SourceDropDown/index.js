import React from "react";
import PropTypes from "prop-types";
import {Select} from "antd";

const {Option} = Select;

export default class SourceDropDown extends React.Component {
    static propTypes = {
        onChange: PropTypes.func,
        provider: PropTypes.object,
    };

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

    onChange = (sourceId) => {
        sourceId = parseInt(sourceId);

        this.setState({value: sourceId});

        if (this.props.onChange)
            this.props.onChange(sourceId);
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
                        <Option key={`source-${index}`}
                                value={object.id}>{`${object.title} [${object.kind}] ${object.address}`}</Option>
                    )
                })}
            </Select>
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
