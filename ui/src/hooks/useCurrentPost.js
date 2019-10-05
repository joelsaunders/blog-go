import {useEffect, useState} from 'react';
import theBookOfJoel from "../apis/theBookOfJoel";

export default (postSlug) => {
    const [currentPost, setCurrentPost] = useState(null);

    useEffect(() => {
        (async () => {
            const response = await theBookOfJoel.get(
                `/api/v1/posts/${postSlug}`
            );
            setCurrentPost(response.data);
        })();
    }, [postSlug]);

    return currentPost
};