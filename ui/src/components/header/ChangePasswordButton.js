import React from 'react';

export const ChangePasswordButton = (props) => {
    return <button
        onClick={props.handleToggleOpen}
        className="inline-block text-sm px-4 py-2 leading-none border rounded text-white border-white hover:border-transparent hover:text-teal-500 hover:bg-white mt-4 md:my-auto lg:mt-0"
    >
        Change Password
    </button>;
};