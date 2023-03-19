import React from "react";

import {shallow} from "enzyme";

import {Select} from "antd";
import SourceDropDown from "./index";

const {Option} = Select;

describe('<SourceDropDown /> rendering', () => {
    it('renders correctly', async () => {
        const client = new MockDashboardMainApiClient();
        const wrapper = shallow(<SourceDropDown provider={client}/>);

        expect(wrapper.find(Select).find(Option)).toHaveLength(1); // only "no data" option

        await fakeSourcesListPromise;
        wrapper.update();

        expect(wrapper.find(Select).find(Option)).toHaveLength(4); // 4 option - "none" and 3 sources
        expect(wrapper).toMatchSnapshot();
    });

    it('onChange callback works', async () => {
        const expectedValue = 101;
        let actualValue = 0;
        let onChangeSpy = function (sourceId) {
            actualValue = sourceId;
        }

        const wrapper = shallow(<SourceDropDown provider={new MockDashboardMainApiClient()} onChange={onChangeSpy}/>);
        expect(wrapper.find(Select)).toHaveLength(1);
        expect(wrapper.find(Select).find(Option)).toHaveLength(1); // only "no data" option

        await fakeSourcesListPromise;
        wrapper.update();

        expect(wrapper.find(Select)).toHaveLength(1);
        expect(wrapper.find(Select).find(Option)).toHaveLength(4); // 4 option - "none" and 3 sources

        wrapper.find(Select).simulate("change", "101");
        expect(actualValue).toEqual(expectedValue);
    });
});


const fakeSourcesListPromise = new Promise((resolve, reject) => {
    const fakeSourcesList = [
        {id: 101, title: "Source 1", kind: "prometheus", address: "http://foo-bar:4001/"},
        {id: 102, title: "Source 2", kind: "influxdb", address: "http://foo-bar:4002/"},
        {id: 103, title: "Source 3", kind: "prometheus", address: "http://foo-bar:4003/"},
    ];

    resolve(fakeSourcesList)
});

class MockDashboardMainApiClient {
    loadSourcesList() {
        return fakeSourcesListPromise;
    }
}