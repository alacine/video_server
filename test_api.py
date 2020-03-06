#!/usr/bin/python3

import json
import requests


user_url = 'http://localhost:8000/api/users/'
video_url = 'http://localhost:8000/api/videos/'
header = {
    'X-User-Name': 'test_man'
}


def test_users():
    payload = {
        'user_name': 'test_man',
        'pwd': 'test_password'
    }
    payload = json.dumps(payload)
    uid = '6'
    # create user
    print('create user')
    response = requests.post(user_url, data=payload)
    print(response.json())
    # get user
    print('get user')
    response = requests.get(user_url+uid, headers=header)
    print(response.json())
    # list user videos
    print('list user videos')
    response = requests.get(user_url+uid+'/videos')
    print(response.json())


def test_videos():
    payload = {
        'author_id': 6,
        'title': 'first_video',
        'description': 'first_video_description'
    }
    payload = json.dumps(payload)
    vid = '2'
    # get all videos
    print('get all videos')
    response = requests.get(video_url)
    print(response.json())
    # get single video info
    print('get single video info')
    response = requests.get(video_url+vid)
    print(response.json())
    # add new video
    print('add new video')
    response = requests.post(video_url, data=payload, headers=header)
    print(response.json())
    # delete video
    print('delete video')
    response = requests.delete(video_url+vid, headers=header)
    print(response.text)

    payload = {
        'author_id': 6,
        'content': 'first_video_first_comment'
    }
    payload = json.dumps(payload)
    cm = '/comments'
    response = requests.post(video_url+vid+cm, data=payload, headers=header)
    # post comment
    print('post comment')
    print(response.text)
    # get comments
    print('get comments')
    response = requests.get(video_url+vid+cm)
    print(response.json())

def main():
    test_users()
    test_videos()


if __name__ == "__main__":
    main()
