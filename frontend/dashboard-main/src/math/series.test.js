import {getAverageDiff} from "./series";

describe('getAverageDiff', () => {
    it('no changes', async () => {
        const series1 = [100, 100, 100, 100];
        const series2 = [100, 100, 100, 100];
        const target = getAverageDiff(series1, series2);

        expect(target.absolute).toEqual(0);
        expect(target.percentage).toEqual(0);
    });

    it('basic', async () => {
        const series1 = [100, 100, 100, 100];
        const series2 = [90, 90, 90, 90];
        const target = getAverageDiff(series1, series2);

        expect(target.absolute).toEqual(-10);
        expect(target.percentage).toEqual(-10);
    });

    it('basic 2', async () => {
        const series1 = [100, 100, 100, 100];
        const series2 = [10, 90, 10, 90];
        const target = getAverageDiff(series1, series2);

        expect(target.absolute).toEqual(-50);
        expect(target.percentage).toEqual(-50);
    });

    it('no source data and increase', async () => {
        const series1 = [0, 0, 0, 0];
        const series2 = [10, 10, 10, 10];
        const target = getAverageDiff(series1, series2);

        expect(target.absolute).toEqual(10);
        expect(target.percentage).toEqual(Infinity);
    });

    it('no source data and decrease', async () => {
        const series1 = [0, 0, 0, 0];
        const series2 = [-10, -10, -10, -10];
        const target = getAverageDiff(series1, series2);

        expect(target.absolute).toEqual(-10);
        expect(target.percentage).toEqual(-Infinity);
    });

    it('no source data and no current data', async () => {
        const series1 = [0, 0, 0, 0];
        const series2 = [0, 0, 0, 0];
        const target = getAverageDiff(series1, series2);

        expect(target.absolute).toEqual(0);
        expect(target.percentage).toEqual(0);
    });
});