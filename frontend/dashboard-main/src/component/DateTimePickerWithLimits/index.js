import React from "react";
import DateTime from "react-datetime";
import "react-datetime/css/react-datetime.css";
// import {Moment} from 'moment';
import PropTypes from "prop-types";

class DateTimePickerWithLimits extends React.Component {
    static propTypes = {
        from: PropTypes.number,
        till: PropTypes.number,
        onChange: PropTypes.func,
    };

    constructor(props) {
        super(props);
        this.state = {
            value: new Date(parseInt(this.props.from) * 1000),
        }
    }

    componentDidUpdate(prevProps) {
        if (prevProps.from !== this.props.from) {
            this.setState({
                value: new Date(parseInt(this.props.from) * 1000),
            })
        }
    }

    render() {
        return (
            <div>
                <DateTime value={this.state.value}
                          utc={true}
                          onChange={this.onChange}
                />
                <span>{new Date(this.props.from * 1000).toUTCString()} - {new Date(this.props.till * 1000).toUTCString()}</span>
            </div>
        )
    }

    /**
     *
     * @param {Moment} moment
     */
    onChange = (moment) => {
        this.setState({value: moment.toDate()})

        if (this.props.onChange) {
            this.props.onChange(moment.unix());
        }
    }
}

export default DateTimePickerWithLimits;