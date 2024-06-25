const DateFormater = (date) => {
    const d = new Date(date);
    return d.toUTCString();
}

export default DateFormater;

