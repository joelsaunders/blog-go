import React from 'react';
import Modal from "../Modal";
import LoginForm from "../forms/LoginForm";
import {Link} from "react-router-dom";
import {LogInButton} from "./LogInButton";
import HeaderMenu from "./HeaderMenu";
import {ChangePasswordButton} from "./ChangePasswordButton";
import ChangePasswordForm from "../forms/ChangePasswordForm";


const renderLoginModal = (onDismissLogIn) => {
    return <Modal onDismiss={onDismissLogIn}>
        <LoginForm onDismiss={onDismissLogIn} />
    </Modal>
};

const renderPasswordChange = (toggleOpen) => {
    return <Modal onDismiss={toggleOpen}>
        <ChangePasswordForm onDismiss={toggleOpen} />
    </Modal>
};

const Header = (props) => {
    const {menuOpen, onMenuClick} = props;

    return <nav className="flex items-center justify-between flex-wrap bg-teal-500 px-6 py-3">
        <div className="flex items-center flex-shrink-0 text-white mr-6">
            <Link to="/">
                <svg xmlns="http://www.w3.org/2000/svg" xmlnsXlink="http://www.w3.org/1999/xlink"
                     enableBackground="new 0 0 512 512" id="Layer_1" version="1.1"
                     viewBox="0 0 512 512" width="56" height="56"
                     strokeLinecap="round" strokeLinejoin="round"
                     stroke="currentColor" strokeWidth="6" fill="currentColor"
                     className="text-white pr-2"
                >
                    <g><g><g><path d="M501.3,415.9c9.7,0,9.7-15,0-15C491.7,400.9,491.7,415.9,501.3,415.9L501.3,415.9z"/></g></g><g><g><path d="M114.5,415.9c21.5,0,43.1,0,64.6,0c35.2,0,70.4,0,105.5,0c13.2,0,26.3,0,39.5,0c32.4,0,64.8,0,97.2,0     c26.2,0,52.5,0.7,78.7,0c0.4,0,0.9,0,1.3,0c5.5,0,9.5-6.4,6.5-11.3c-8.6-13.7-17.2-27.4-25.8-41.1     c-17.7-28.2-35.5-56.4-53.2-84.7c-15.1-24-30.2-48-45.3-72c-6.7-10.7-13.5-21.4-20.2-32.1c-15.4-24.4-30.5-49-46.1-73.3     c-0.3-0.5-0.7-1.1-1-1.6c-2.8-4.5-10.3-5.3-13,0c-3.7,7.3-8.9,14.2-13.3,21.2c-8.9,14.2-17.8,28.3-26.7,42.5     c-10.7,17-21.4,34.1-32.1,51.1c-7.4,11.8-14.9,23.6-22.3,35.4c-5.2,8.2,7.8,15.7,13,7.6c11.8-18.7,23.6-37.4,35.3-56.2     c14.3-22.7,28.5-45.4,42.8-68.1c5.4-8.5,11.8-17,16.3-26c-4.3,0-8.6,0-13,0c13.7,21.8,27.4,43.6,41.1,65.3     c7.7,12.2,15.3,24.4,23,36.6c13.9,22.2,27.9,44.3,41.8,66.5c18.1,28.7,36.1,57.4,54.2,86.1c6.9,11,13.8,22,20.7,33     c3.4,5.4,6.5,11.5,10.4,16.6c0.2,0.2,0.3,0.5,0.5,0.7c2.2-3.8,4.3-7.5,6.5-11.3c-22.1,0-44.3,0-66.4,0c-34.2,0-68.5,0-102.7,0     c-12.7,0-25.5,0-38.2,0c-33.8,0-67.6,0-101.4,0c-25.6,0-51.3-0.7-76.9,0c-0.4,0-0.9,0-1.3,0C104.8,400.9,104.8,415.9,114.5,415.9     L114.5,415.9z"/></g></g><g><g><path d="M359.7,197.2c0.2,0.4,0.3,0.7,0.5,1.1c-0.3-1.3-0.7-2.5-1-3.8c-0.1-1.1,2.7-1.7,0.9-1.5c-0.8,0.1-1.7,0.8-2.5,1.2     c-1.8,0.8-3.6,1.7-5.4,2.5c-2.8,1.3-5.5,2.6-8.3,3.9c-1.2,0.5-3.9,1.2-4.7,2.2c-0.4,0.5-3.5,0.2,1.6,0.2c4.4,0,2,0.2,1.4-0.3     c-2.1-2-6.5-3.1-9.1-4.3c-6.5-3-12.9-6.1-19.4-9.1c-3.7-1.7-6.3-1.4-9.8,0.2c-3.1,1.4-6.2,2.9-9.2,4.4c-7,3.3-14.1,6.6-21.1,10     c2.5,0,5,0,7.6,0c-4.2-1.9-8.4-3.9-12.6-5.8c-1.3-0.6-2.6-1.2-3.9-1.8c-1.4-0.7-2.9-1.3-4.3-2c-0.4-0.2-1.1-0.7-1.5-0.7     c-1.5,0,1.7-2.9,1.1,1.6c-0.3,1.3-0.7,2.5-1,3.8c1.5-3.5,4-6.6,6-9.8c4.5-7.2,9-14.4,13.5-21.5c10.3-16.4,20.6-32.8,30.9-49.2     c2.3-3.7,4.7-7.3,6.9-11c-4.3,0-8.6,0-13,0c16.5,26.2,32.9,52.4,49.4,78.5C354.9,189.7,357.3,193.4,359.7,197.2     c5.1,8.2,18.1,0.6,13-7.6c-16.5-26.2-32.9-52.4-49.4-78.5c-2.4-3.8-4.7-7.5-7.1-11.3c-3-4.7-10-5-13,0     c-12.9,22-27.2,43.3-40.8,64.8c-5.3,8.4-11.3,16.7-16,25.4c-2.4,4.4-2.3,9.7,1.1,13.5c2.5,2.9,6.5,4.2,9.9,5.8     c5.1,2.4,10.2,4.8,15.4,7.2c4,1.8,6.8,0.9,10.4-0.7c3.5-1.7,7-3.3,10.5-5c6.3-3,12.7-6,19-8.9c-2.5,0-5,0-7.6,0     c5.3,2.5,10.5,5,15.8,7.4c3.6,1.7,7.2,3.4,10.7,5.1c2.6,1.2,5.4,2.8,8.4,3c6.3,0.3,13.7-4.9,19.2-7.5c2.8-1.3,5.8-2.5,8.6-4     c6.3-3.5,7.6-10,4.8-16.2c-1.7-3.7-7-4.5-10.3-2.7C358.5,189.1,358.1,193.5,359.7,197.2z"/></g></g><g><g><path d="M312.5,400.9c-34.1,0-68.1,0-102.2,0c-54.2,0-108.4,0-162.6,0c-12.3,0-24.6,0-36.9,0c2.2,3.8,4.3,7.5,6.5,11.3     c20.7-33,41.5-66,62.2-98.9c17.2-27.4,34.4-54.7,51.6-82.1c12.4-19.7,24.7-39.3,37.1-59c-4.3,0-8.6,0-13,0     c19.8,31.5,39.6,63,59.5,94.6c21.4,34,42.7,67.9,64.1,101.9c9.1,14.5,18.3,29,27.4,43.6c5.1,8.2,18.1,0.6,13-7.6     c-21.2-33.7-42.3-67.3-63.5-101c-16.7-26.5-33.4-53.1-50.1-79.6c-12.4-19.8-24.9-39.6-37.3-59.3c-3-4.8-9.9-4.8-13,0     c-16.9,26.9-33.8,53.8-50.7,80.7c-25.1,39.9-50.1,79.7-75.2,119.6c-8.3,13.2-16.6,26.5-25,39.7c-3.1,4.9,0.9,11.3,6.5,11.3     c34.1,0,68.1,0,102.2,0c54.2,0,108.4,0,162.6,0c12.3,0,24.6,0,36.9,0C322.1,415.9,322.1,400.9,312.5,400.9z"/><path d="M312.5,415.9c9.7,0,9.7-15,0-15C302.8,400.9,302.8,415.9,312.5,415.9L312.5,415.9z"/></g></g><g><g><path d="M204.9,236.8c-7.6,3.6-15.2,7.1-22.9,10.7c2.5,0,5,0,7.6,0c-3.9-1.8-7.8-3.7-11.7-5.5c-4.5-2.1-12.1-7.5-17.3-6.8     c-5,0.7-10.5,4.5-14.9,6.6c-4,1.9-8,3.8-12,5.7c2.5,0,5,0,7.6,0c-7.7-3.6-15.3-7.1-23-10.7c0.9,3.4,1.8,6.8,2.7,10.3     c13.7-21.8,27.4-43.6,41.2-65.5c2-3.1,3.9-6.3,5.9-9.4c-4.3,0-8.6,0-13,0c13.7,21.8,27.4,43.6,41.2,65.5c2,3.1,3.9,6.3,5.9,9.4     c5.1,8.2,18.1,0.6,13-7.6c-13.7-21.8-27.4-43.6-41.2-65.5c-2-3.1-3.9-6.3-5.9-9.4c-3-4.8-9.9-4.8-13,0     c-13.7,21.8-27.4,43.6-41.2,65.5c-2,3.1-3.9,6.3-5.9,9.4c-2,3.2-0.9,8.6,2.7,10.3c4.3,2,8.7,4,13,6.1c4.7,2.2,11.4,6.9,16.6,5     c8.5-3,17-7.7,25.1-11.8c-2.5,0-5,0-7.6,0c4.7,2.2,9.4,4.4,14.1,6.6c3.8,1.8,7.8,4.4,11.9,5.5c4.7,1.3,11.4-3.4,15.6-5.3     c4.4-2,8.7-4.1,13.1-6.1C221.2,245.7,213.6,232.8,204.9,236.8z"/><path d="M208.6,250.8c9.7,0,9.7-15,0-15C199,235.8,199,250.8,208.6,250.8L208.6,250.8z"/></g></g></g>
                </svg>
            </Link>
            <Link to="/">
                <span className="font-semibold text-xl tracking-tight">Joel Saunders</span>
            </Link>
        </div>
        <div className="block md:hidden">
            <button
                className="flex items-center px-3 py-2 border rounded text-teal-200 border-teal-400 hover:text-white hover:border-white"
                onClick={onMenuClick}
            >
                <svg className="fill-current h-3 w-3" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg">
                    <title>Menu</title>
                    <path d="M0 3h20v2H0V3zm0 6h20v2H0V9zm0 6h20v2H0v-2z"/>
                </svg>
            </button>
        </div>
        <div className={`w-full block flex-grow md:flex md:items-center md:w-auto ${menuOpen? null: 'hidden'}`}>
            <HeaderMenu currentPath={props.currentPath} />
        </div>
        <div className={menuOpen? '': 'hidden md:block'}>
            {
                props.loggedIn?
                    <ChangePasswordButton handleToggleOpen={props.handleChangePasswordClick} />:
                    null
            }
            <div className="ml-2 inline-block">
                <LogInButton
                    loggedIn={props.loggedIn}
                    handleLogInClick={props.handleLogInClick}
                    handleLogOutClick={props.handleLogOutClick}
                />
            </div>
        </div>
        {props.renderLogInModal? renderLoginModal(props.onDismissLogIn) : null}
        {props.passwordChangeOpen? renderPasswordChange(props.handleChangePasswordClick): null}
    </nav>;
};

export default Header
