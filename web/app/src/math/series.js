function getAverageDiff(seriesFrom, seriesTo) {
    let sumOne = 0.0;
    let sumTwo = 0.0;
    seriesFrom.forEach((val) => sumOne += val);
    seriesTo.forEach((val) => sumTwo += val);

    const valOne = sumOne / seriesFrom.length;
    const valTwo = sumTwo / seriesTo.length;

    let percentage = -1 * (100 - valTwo / (valOne / 100));

    if (percentage === 0) percentage = 0; // to prevent negative zeros
    if (isNaN(percentage)) percentage = 0; // to prevent negative zeros

    return {
        absolute: valTwo - valOne,
        percentage: percentage,
    }
}

module.exports = {
    getAverageDiff: getAverageDiff,
}