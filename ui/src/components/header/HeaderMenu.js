import React from 'react';
import {Link} from "react-router-dom";
import _ from 'lodash';

const HeaderItems = {
    'Blog': '/',
    // 'About': '/about',
    // 'Site': '/site',
    'Team': '/team',
    'Contact': '/contact',
};

const MenuItem = ({title, path, currentPath}) => {
    const highlighted = currentPath === path ? 'text-teal-100': 'text-teal-200';
    return <Link to={path}>
        <div className={`block mt-4 md:inline-block md:mt-0 hover:text-white mr-4 ${highlighted}`}>
            {title}
        </div>
    </Link>
};

const HeaderMenu = ({currentPath}) => {
    return <div>
        {_.map(HeaderItems, (value, key) => {
            return <MenuItem key={key} title={key} path={value} currentPath={currentPath} />
        })}
    </div>
};

export default HeaderMenu;
