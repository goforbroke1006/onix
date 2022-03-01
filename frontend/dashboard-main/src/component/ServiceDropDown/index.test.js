import React from "react";
import {mount, shallow} from "enzyme";

import ServiceDropDown from "./index";

describe('<ServiceDropDown /> rendering', () => {
    it('renders correctly', () => {
        const client = new MockDashboardMainApiClient();
        const wrapper = mount(<ServiceDropDown provider={client}/>);

        fakeServiceListPromise.then(() => {
            wrapper.update();
            expect(wrapper).toMatchSnapshot();
        });
    });

    it('onChange callback works', () => {
        const expectedValue = "foo/backend";
        let actualValue = null;

        let onChangeSpy = function (serviceName) {
            actualValue = serviceName;
        }

        const wrapper = shallow(<ServiceDropDown provider={new MockDashboardMainApiClient()} onChange={onChangeSpy}/>);
        wrapper.find('select').at(0).simulate('change', {
            target: {value: 'foo/backend', name: 'foo/backend'}
        });

        expect(actualValue).toEqual(expectedValue);
    });
});

const fakeServiceListPromise = new Promise((resolve, reject) => {
    const fakeServicesList = [
        {title: "acme/backend"},
        {title: "foo/backend"},
        {title: "bar/backend"},
    ];

    resolve(fakeServicesList)
});

class MockDashboardMainApiClient {
    loadServices() {
        return fakeServiceListPromise;
    }
}