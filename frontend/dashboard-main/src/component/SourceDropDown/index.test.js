import React from "react";
import {mount, shallow} from "enzyme";

import SourceDropDown from "./index";

describe('<SourceDropDown /> rendering', () => {
    it('renders correctly', () => {
        const client = new MockDashboardMainApiClient();
        const wrapper = mount(<SourceDropDown provider={client}/>);

        fakeSourcesListPromise.then(() => {
            wrapper.update();
            expect(wrapper).toMatchSnapshot();
        });
    });

    it('onChange callback works', () => {
        const expectedValue = 101;
        let actualValue = 0;
        let onChangeSpy = function (sourceId) {
            actualValue = sourceId;
        }

        const wrapper = shallow(<SourceDropDown provider={new MockDashboardMainApiClient()} onChange={onChangeSpy}/>);
        wrapper.find('select').at(0).simulate('change', {
            target: {value: '101', name: 'item1'}
        });

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