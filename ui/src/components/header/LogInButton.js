import React from "react";

export const LogInButton = ({loggedIn, handleLogInClick, handleLogOutClick}) => {
    return <button
        onClick={loggedIn? handleLogOutClick: handleLogInClick}
        className="inline-block text-sm px-4 py-2 leading-none border rounded text-white border-white hover:border-transparent hover:text-teal-500 hover:bg-white mt-4 md:my-auto lg:mt-0"
    >
        {loggedIn? 'Log Out': 'Log In'}
    </button>;
};