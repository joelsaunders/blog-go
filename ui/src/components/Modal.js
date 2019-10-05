import React from 'react';
import ReactDOM from 'react-dom';

const Modal = (props) => {
    return ReactDOM.createPortal(
        <div onClick={props.onDismiss} className="absolute top-0 flex w-full h-screen justify-center bg-smoke-dark">
            <div onClick={(event) => event.stopPropagation()} className="my-auto">
                <div className="bg-white shadow-md rounded px-8 pt-6 pb-8 mb-4100">{props.children}</div>
            </div>
        </div>,
        document.querySelector('#modal')
    )
};

export default Modal
