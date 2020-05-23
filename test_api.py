#!/usr/bin/python3

import json
import requests


user_url = 'http://localhost:8000/api/users/'
session_url = 'http://localhost:8000/api/sessions/'
video_url = 'http://localhost:8000/api/videos/'
cookie = {
    'X-User-Id': '',
    'X-Session-Id': ''
}


def test_users():
    payload = {
        'user_name': 'test_man_py',
        'pwd': 'test_password_py'
    }
    payload = json.dumps(payload)
    # create user
    print('create user')
    response = requests.post(user_url, data=payload)
    uid = str(response.json()['user_id'])
    print('   ', response.json())
    # get user
    print('get user')
    response = requests.get(user_url+uid)
    print('   ', response.json())
    # list user videos
    print('list user videos')
    response = requests.get(user_url+uid+'/videos')
    print('   ', response.json())
    # user login
    print('login')
    response = requests.post(session_url, data=payload)
    cookie['X-User-Id'] = str(response.json()['user_id'])
    cookie['X-Session-Id'] = response.json()['session_id']
    print('   ', response.json())


def test_videos():
    payload = {
        'author_id': int(cookie['X-User-Id']),
        'title': 'first_video',
        'description': 'first_video_description'
    }
    payload = json.dumps(payload)
    vid = '1'
    # get all videos
    print('get all videos')
    response = requests.get(video_url)
    print('   ', response.json())
    # get single video info
    print('get single video info')
    response = requests.get(video_url+vid)
    print('   ', response.json())
    # add new video
    print('add new video')
    response = requests.post(video_url, data=payload, cookies=cookie)
    print('   ', response.json())
    # delete video
    print('delete video')
    response = requests.delete(video_url+vid, cookies=cookie)
    print('   ', response.text)

    payload = {
        'author_id': 6,
        'content': 'first_video_first_comment'
    }
    payload = json.dumps(payload)
    cm = '/comments'
    response = requests.post(video_url+vid+cm, data=payload, cookies=cookie)
    # post comment
    print('post comment')
    print('   ', response.text)
    # get comments
    print('get comments')
    response = requests.get(video_url+vid+cm)
    print('   ', response.json())


def logout():
    print('logout')
    response = requests.delete(session_url, cookies=cookie)
    print('   ', response.text)


def main():
    test_users()
    test_videos()
    logout()


if __name__ == "__main__":
    main()
