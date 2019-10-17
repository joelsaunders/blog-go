import React from 'react';
import {connect} from 'react-redux';
import Modal from "../../Modal";
import {deletePost} from "../../../actions";


const DeletePostModal = ({setDeleteModalActive, postSlug, deletePost}) => {
    return <Modal onDismiss={() => {setDeleteModalActive(false)}}>
        Are you sure you want to delete this post?
        <button onClick={() => {deletePost(postSlug)}}>Delete</button>
        <button onClick={() => {setDeleteModalActive(false)}}>Cancel</button>
    </Modal>;
};

export default connect(null, {deletePost})(DeletePostModal);
