import React, {useState, useEffect} from 'react';
import {connect} from "react-redux";
import {signOut} from "../../actions";
import Header from "./Header";


const HeaderContainer = (props) => {
    const[logInOpen, setLoginOpen] = useState(false);
    const[menuOpen, setMenuOpen] = useState(false);

    const onMenuClick = () => {
        setMenuOpen(!menuOpen);
    };

    const onLogInClick = () => {
        setLoginOpen(true);
    };

    const onLogOutClick = () => {
        props.signOut()
    };

    const onDismissLogIn = () => {
        setLoginOpen(false);
    };

    return <Header
            renderLogInModal={logInOpen}
            handleLogInClick={onLogInClick}
            handleLogOutClick={onLogOutClick}
            loggedIn={props.loggedIn}
            onDismissLogIn={onDismissLogIn}
            onMenuClick={onMenuClick}
            menuOpen={menuOpen}
            currentPath={props.location.pathname}
    />
};

const mapStateToProps = (state) => {
    return {loggedIn: state.auth.token !== null}
};

export default connect(mapStateToProps, {signOut})(HeaderContainer);
