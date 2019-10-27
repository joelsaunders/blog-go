import _ from 'lodash';
import React, {useEffect, useState} from 'react';
import Modal from "../Modal";
import {connect} from "react-redux";
import {fetchTags, editTag, createTag, deleteTag} from "../../actions";
import {DebounceInput} from "react-debounce-input";


const TagItem = ({tag, editTag, deleteTag}) => {
    const [inputValue, setInputValue] = useState(tag.name);

    return <div className="my-1">
        <DebounceInput
            value={inputValue}
            debounceTimeout={500}
            onChange={event => {
                setInputValue(event.target.value);
                editTag(tag.id, event.target.value);
            }}
        />
        <button
            onClick={() => deleteTag(tag.id)}
            className="rounded text-red-500 text-sm border-red-500 border px-2 hover:text-red-100 hover:bg-red-500"
        >
            Delete
        </button>
    </div>;
};

const TagModal = ({tags, fetchTags, editTag, createTag, deleteTag, onDismiss}) => {
    const [newTagValue, setNewTagValue] = useState("");

    useEffect(
        () => {
            (async () =>  await fetchTags())()
        },
        [fetchTags]
    );

    return <Modal onDismiss={onDismiss}>
        Manage Tags
        {_.map(tags, (tag) => <TagItem key={tag.id} tag={tag} editTag={editTag} deleteTag={deleteTag} />)}
        <input
            value={newTagValue}
            placeholder="enter new tag name"
            onChange={event => {setNewTagValue(event.target.value)}}
        />
        <button
            className="mr-2 rounded text-teal-500 text-sm border-teal-500 border px-2 hover:text-teal-100 hover:bg-teal-500"
            onClick={() => {
                createTag(newTagValue);
                setNewTagValue("");
            }}
        >Submit</button>
    </Modal>;
};

const mapStateToProps = (state) => {
    return {tags: state.tags}
};

export default connect(mapStateToProps, {fetchTags, editTag, createTag, deleteTag})(TagModal);
