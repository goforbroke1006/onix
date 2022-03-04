import React from "react";
// import DateTime from "react-datetime";
import "react-datetime/css/react-datetime.css";
// import moment from 'moment';
import PropTypes from "prop-types";
import {DatePicker} from "antd";
import moment from 'moment-timezone';

// import locale from 'antd/es/date-picker/locale/';

class DateTimePickerWithLimits extends React.Component {
    static propTypes = {
        from: PropTypes.number,
        till: PropTypes.number,
        onChange: PropTypes.func,
    };

    constructor(props) {
        super(props);
        this.state = {
            value: moment.tz(new Date(this.props.from * 1000), "UTC"),
        }
    }

    componentDidUpdate(prevProps) {
        if (prevProps.from !== this.props.from) {
            let m = moment.tz(new Date(this.props.from * 1000), "UTC");
            this.setState({value: m});
            if (this.props.onChange) {
                this.props.onChange(m.unix());
            }
        }
    }

    render() {
        return (
            <div>
                <DatePicker value={moment(this.state.value, 'YYYY-MM-DD')} showTime={true} onChange={this.onChange}/>
                {/*<DateTime value={this.state.value}*/}
                {/*          utc={true}*/}
                {/*          onChange={this.onChange}*/}
                {/*/>*/}
                <br/>
                <span>{new Date(this.props.from * 1000).toUTCString()} - {new Date(this.props.till * 1000).toUTCString()}</span>
            </div>
        )
    }

    /**
     *
     * @param {Moment} date
     * @param {String} dateString
     */
    onChange = (date, dateString) => {
        console.log(dateString);

        this.setState({value: date});

        if (this.props.onChange) {
            this.props.onChange(date.unix());
        }
    }
}

export default DateTimePickerWithLimits;