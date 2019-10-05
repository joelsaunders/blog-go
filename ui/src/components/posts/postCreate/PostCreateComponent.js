import React from 'react';


const PostCreateComponent = ({children}) => {
    return <div className="w-full my-6 rounded overflow-hidden shadow-lg bg-white">
        <div className="px-6 py-4">
            <h3 className="text-gray-900 font-bold text-xl mb-2">
                Create a Post
            </h3>
            {children}
        </div>
    </div>;
};

export default PostCreateComponent;
