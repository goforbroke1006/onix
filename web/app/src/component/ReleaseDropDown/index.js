import React from "react";
import PropTypes from "prop-types";

import {Select} from "antd";

const {Option} = Select;

const defaultValue = 0;

class ReleaseDropDown extends React.Component {
    static propTypes = {
        serviceName: PropTypes.string,
        onChange: PropTypes.func,
        provider: PropTypes.object,
    };

    constructor(props) {
        super(props);
        this.state = {
            value: defaultValue,
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
    }

    render() {
        let emptyOption;
        if (this.state.items.length === 0) {
            emptyOption = (<Option value={defaultValue}>-- no data --</Option>)
        } else {
            emptyOption = (<Option value={defaultValue}>-- none --</Option>)
        }

        return (
            <Select defaultValue={this.state.value} style={{width: 600}} onChange={this.onChange}>
                {emptyOption}
                {this.state.items.map(function (object) {
                    return (
                        <Option key={"service-select-option-" + object.id}
                                value={object.id}>{object.title}</Option>
                    )
                })}
            </Select>
        )
    }

    onChange = (releaseId) => {
        releaseId = parseInt(releaseId);

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

        this.props.provider.loadReleasesList(this.props.serviceName)
            .then(data => {
                console.debug("load releases");
                this.setState({items: data})
            });
    }
}

export default ReleaseDropDown;