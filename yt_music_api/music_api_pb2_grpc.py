# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
"""Client and server classes corresponding to protobuf-defined services."""
import grpc

import music_api_pb2 as music__api__pb2


class MusicApiStub(object):
    """Missing associated documentation comment in .proto file."""

    def __init__(self, channel):
        """Constructor.

        Args:
            channel: A grpc.Channel.
        """
        self.SearchArtists = channel.unary_unary(
                '/music_api.MusicApi/SearchArtists',
                request_serializer=music__api__pb2.Query.SerializeToString,
                response_deserializer=music__api__pb2.Artists.FromString,
                )
        self.SearchAlbums = channel.unary_unary(
                '/music_api.MusicApi/SearchAlbums',
                request_serializer=music__api__pb2.Query.SerializeToString,
                response_deserializer=music__api__pb2.Albums.FromString,
                )


class MusicApiServicer(object):
    """Missing associated documentation comment in .proto file."""

    def SearchArtists(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def SearchAlbums(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')


def add_MusicApiServicer_to_server(servicer, server):
    rpc_method_handlers = {
            'SearchArtists': grpc.unary_unary_rpc_method_handler(
                    servicer.SearchArtists,
                    request_deserializer=music__api__pb2.Query.FromString,
                    response_serializer=music__api__pb2.Artists.SerializeToString,
            ),
            'SearchAlbums': grpc.unary_unary_rpc_method_handler(
                    servicer.SearchAlbums,
                    request_deserializer=music__api__pb2.Query.FromString,
                    response_serializer=music__api__pb2.Albums.SerializeToString,
            ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
            'music_api.MusicApi', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))


 # This class is part of an EXPERIMENTAL API.
class MusicApi(object):
    """Missing associated documentation comment in .proto file."""

    @staticmethod
    def SearchArtists(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/music_api.MusicApi/SearchArtists',
            music__api__pb2.Query.SerializeToString,
            music__api__pb2.Artists.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def SearchAlbums(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/music_api.MusicApi/SearchAlbums',
            music__api__pb2.Query.SerializeToString,
            music__api__pb2.Albums.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)