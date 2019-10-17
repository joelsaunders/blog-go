
export default function postDateFormatter(dateString) {
    return new Intl.DateTimeFormat('en-GB', {
        year: 'numeric',
        month: 'long',
        day: '2-digit'
    }).format(Date.parse(dateString))
}
