import React from "react";

import {shallow} from "enzyme";

import {Select} from "antd";
const {Option} = Select;
import ServiceDropDown from "./index";

describe('<ServiceDropDown /> component', () => {
    it('renders correctly', async () => {
        const client = new MockDashboardMainApiClient();
        const wrapper = shallow(<ServiceDropDown provider={client}/>);

        expect(wrapper.find(Select).find(Option)).toHaveLength(1); // only "no data" option

        await fakeServiceListPromise;
        wrapper.update();

        expect(wrapper.find(Select).find(Option)).toHaveLength(4); // 4 option - "none" and 3 services
        expect(wrapper).toMatchSnapshot();
    });

    it('onChange callback works', async () => {
        const expectedValue = "foo/backend";
        let actualValue = null;

        let onChangeSpy = function (serviceName) {
            actualValue = serviceName;
        }

        const wrapper = shallow(<ServiceDropDown provider={new MockDashboardMainApiClient()} onChange={onChangeSpy}/>);
        expect(wrapper.find(Select)).toHaveLength(1);
        expect(wrapper.find(Select).find(Option)).toHaveLength(1); // only "no data" option

        await fakeServiceListPromise;
        wrapper.update();

        expect(wrapper.find(Select)).toHaveLength(1);
        expect(wrapper.find(Select).find(Option)).toHaveLength(4); // 4 option - "none" and 3 services

        wrapper.find(Select).simulate("change", "foo/backend");
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