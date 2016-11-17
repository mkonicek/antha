// Copyright 2011 <chaishushan@gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// OpenCV Homepage: http://code.opencv.org

#ifndef _GO_OPENCV_BINDING_H_
#define _GO_OPENCV_BINDING_H_

#include <opencv/cv.h>
#include <opencv/highgui.h>
#if (CV_VERSION_MAJOR * 100000 + CV_VERSION_MINOR * 100) >= 200400
# include <opencv2/photo/photo_c.h>
#endif
#include <opencv2/imgproc/imgproc_c.h>
#include <opencv2/core/core_c.h>
#include <opencv2/core/version.hpp>

// Trackbar
int GoOpenCV_CreateTrackbar(
	char* trackbar_name, char* window_name,
	int value, int count
);
void GoOpenCV_DestroyTrackbar(
	char* trackbar_name, char* window_name
);

// mouse callback
void GoOpenCV_SetMouseCallback(
	const char* window_name
);

// video writer args
unsigned GoOpenCV_FOURCC_(
	int c1, int c2, int c3, int c4
);

#endif // _GO_OPENCV_BINDING_H_
